//
//  WHOMessage.h
//  Hooli
//
//  Created by dylan on 4/19/14.
//  Copyright (c) 2014 whoisdylan. All rights reserved.
//

#import <Foundation/Foundation.h>
#import <CoreLocation/CoreLocation.h>

@interface WHOMessage : NSObject

- (instancetype) initWithMessage:(NSString* )message Author:(NSString* )author Location:(CLLocation *)location EncodedPhoto:(NSString* )encodedPhoto;
-(BOOL)isEqual:(id)object;
//-(BOOL)isEqualToObject:(WHOMessage* )object;
@property (nonatomic, strong) NSString* message;
@property (nonatomic, strong) NSString* author;
@property (nonatomic, strong) CLLocation* location;
@property (nonatomic, strong) NSString* encodedPhoto;

@end
