//
//  WHOMessage.m
//  Hooli
//
//  Created by dylan on 4/19/14.
//  Copyright (c) 2014 whoisdylan. All rights reserved.
//

#import "WHOMessage.h"

@implementation WHOMessage

-(BOOL)isEqual:(id)object {
    WHOMessage* anObject = (WHOMessage*) object;
    if (self.author == anObject.author) {
//        NSLog(@"not equal 1");
        return NO;
    }
    else if (self.message == anObject.message) {
//        NSLog(@"not equal 2");
        return NO;
    }
    else if (self.location == anObject.location) {
//        NSLog(@"not equal 3");
        return NO;
    }
    else {
//        NSLog(@"message equal");
        return YES;
    }
}

- (instancetype)initWithMessage:(NSString *)message Author:(NSString *)author Location:(CLLocation *)location {
    if (self = [super init]) {
        self.message = message;
        self.author = author;
        self.location = location;
    }
    return self;
}

@end
